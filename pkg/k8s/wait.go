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
	"errors"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cl *Client) WaitForLoadBalancer(namespace, name string, timeoutSeconds int64) (addr string, err error) {
	// Get watch handler to observe changes to services
	watch, err := cl.CoreV1().Services(namespace).Watch(metav1.ListOptions{TimeoutSeconds: &timeoutSeconds})
	if err != nil {
		return
	}

	// Wait for Services to have addresses allocated
	for event := range watch.ResultChan() {
		if event.Type == "Error" || event.Type == "Deleted" {
			err = errors.New("Could not wait for service " + namespace + "/" + name)
			return
		}
		svc, ok := event.Object.(*v1.Service)
		if !ok {
			err = errors.New("Could not wait for service " + namespace + "/" + name)
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
	return
}

func (cl *Client) WaitForPod(namespace, name string, timeoutSeconds int64) error {
	// Get watch handler to observe changes to pods
	watch, err := cl.CoreV1().Pods(namespace).Watch(metav1.ListOptions{TimeoutSeconds: &timeoutSeconds})
	if err != nil {
		return err
	}

	// Wait for pod events
	for event := range watch.ResultChan() {
		if event.Type == "Error" || event.Type == "Deleted" {
			return errors.New("Failed to wait for pod " + namespace + "/" + name)
		}
		// Get the pod
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			return errors.New("Failed to wait for pod " + namespace + "/" + name)
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

		if pod.Status.Phase == "Running" {
			ready := true
			for _, cond := range pod.Status.Conditions {
				if cond.Status != "True" {
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
