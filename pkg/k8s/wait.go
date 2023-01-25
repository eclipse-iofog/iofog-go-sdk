/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package k8s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8swatch "k8s.io/apimachinery/pkg/watch"
)

func getErrorMsg(resource, namespace, name, event string) string {
	return fmt.Sprintf("Failed to wait for %s %s/%s because it is %s", resource, namespace, name, event)
}

func (cl *Client) WaitForLoadBalancer(namespace, name string, timeoutSeconds int64) (addr string, err error) {
	// Get watch handler to observe changes to services
	watch, err := cl.CoreV1().Services(namespace).Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeoutSeconds})
	if err != nil {
		return
	}

	// Wait for Services to have addresses allocated
	for event := range watch.ResultChan() {
		if event.Type == k8swatch.Error || event.Type == k8swatch.Deleted {
			err = errors.New(getErrorMsg("service", namespace, name, string(event.Type)))
			return
		}
		svc, ok := event.Object.(*corev1.Service)
		if !ok {
			err = errors.New(getErrorMsg("service", namespace, name, string(event.Type)))
			return
		}

		// Ignore irrelevant service events
		if svc.Name != name {
			continue
		}
		// Loadbalancer must be ready
		if len(svc.Status.LoadBalancer.Ingress) == 0 {
			continue
		}

		// Check addresses
		ip := svc.Status.LoadBalancer.Ingress[0].IP
		host := svc.Status.LoadBalancer.Ingress[0].Hostname
		if ip != "" {
			addr = ip
		}
		if host != "" {
			addr = host
		}

		if addr == "" {
			continue
		}

		// Return address
		watch.Stop()
	}

	if addr == "" {
		err = errors.New("IP and Hostname values were empty")
	}
	return addr, err
}

func (cl *Client) WaitForPod(namespace, name string, timeoutSeconds int64) error {
	// Get watch handler to observe changes to pods
	watch, err := cl.CoreV1().Pods(namespace).Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeoutSeconds})
	if err != nil {
		return err
	}

	// Wait for pod events
	for event := range watch.ResultChan() {
		if event.Type == k8swatch.Error || event.Type == k8swatch.Deleted {
			return errors.New(getErrorMsg("pod", namespace, name, string(event.Type)))
		}
		// Get the pod
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			return errors.New(getErrorMsg("pod", namespace, name, string(event.Type)))
		}
		// Check pod is in running state
		splitName := strings.Split(pod.Name, "-")
		if len(splitName) <= 2 {
			continue
		}
		splitName = splitName[0 : len(splitName)-2]
		joinName := strings.Join(splitName, "-")
		if joinName != name {
			continue
		}

		if pod.Status.Phase == corev1.PodRunning {
			ready := true
			for _, cond := range pod.Status.Conditions {
				if cond.Status != corev1.ConditionTrue {
					ready = false
					break
				}
			}
			if ready {
				watch.Stop()
			}
		}
	}
	return nil
}
