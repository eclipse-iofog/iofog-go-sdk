package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// readPasswordWithMask reads a password from stdin with asterisk masking
func readPasswordWithMask() (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	var password []byte
	var buf [1]byte

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			if err != nil {
				return "", err
			}
			continue
		}

		char := buf[0]

		// Handle Enter/Return key (13 = carriage return, 10 = line feed)
		if char == 13 || char == 10 {
			// Use \r\n to ensure cursor moves to start of new line
			fmt.Print("\r\n")
			os.Stdout.Sync()
			break
		}

		// Handle backspace (127 = DEL, 8 = backspace)
		if char == 127 || char == 8 {
			if len(password) > 0 {
				password = password[:len(password)-1]
				// Move cursor back, print space to clear asterisk, move cursor back again
				fmt.Print("\b \b")
			}
			continue
		}

		// Handle Ctrl+C
		if char == 3 {
			fmt.Print("\r\n")
			os.Stdout.Sync()
			return "", fmt.Errorf("interrupted by user")
		}

		// Handle Ctrl+D (EOF)
		if char == 4 {
			fmt.Print("\r\n")
			os.Stdout.Sync()
			return "", fmt.Errorf("EOF")
		}

		// Regular character - add to password and echo asterisk
		if char >= 32 && char <= 126 { // Printable ASCII characters
			password = append(password, char)
			fmt.Print("*")
		}
	}

	return string(password), nil
}

// readEmail reads email from stdin
func readEmail() (string, error) {
	fmt.Print("User E-mail: ")
	reader := bufio.NewReader(os.Stdin)
	email, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read email: %v", err)
	}
	return strings.TrimSpace(email), nil
}

// isAuthError checks if the error is an authentication failure
func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Check for HTTP error with 401 or 400 status code (authentication failures)
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.Code == 401 || httpErr.Code == 400
	}
	// Check if error message contains "failed to login"
	return strings.Contains(strings.ToLower(errStr), "failed to login")
}

// create user can be removed!!
func (clt *Client) CreateUser(request User) error {
	// Send request
	if _, err := clt.doRequest("POST", "/user/signup", request); err != nil {
		return err
	}

	return nil
}

func (clt *Client) Login(request LoginRequest) (err error) {
	maxAttempts := 3
	attempt := 0

	for {
		attempt++
		// Prompt for Email if not already set
		if request.Email == "" {
			email, err := readEmail()
			if err != nil {
				return fmt.Errorf("failed to read Email: %v", err)
			}
			request.Email = email
		} else {
			fmt.Printf("User E-mail: %s\n", request.Email)
		}

		// Prompt for Password if not already set
		if request.Password == "" || request.Password == "null" || request.Password == "NULL" {
			fmt.Print("Enter Password: ")
			password, err := readPasswordWithMask()
			if err != nil {
				return fmt.Errorf("failed to read Password: %v", err)
			}
			request.Password = strings.TrimSpace(password)
			// Ensure stdout is flushed before next prompt
			os.Stdout.Sync()
		}

		// Prompt for OTP if not already set
		if request.Totp == "" {
			fmt.Print("Enter OTP: ")
			otp, err := readPasswordWithMask()
			if err != nil {
				return fmt.Errorf("failed to read OTP: %v", err)
			}
			request.Totp = strings.TrimSpace(otp)
		}

		// Send login request
		body, err := clt.doRequest("POST", "/user/login", request)
		if err != nil {
			// Check if it's an authentication error
			if isAuthError(err) {
				if attempt >= maxAttempts {
					return fmt.Errorf("failed to login after %d attempts: %v", maxAttempts, err)
				}
				fmt.Println("\nEmail, password or TOTP is wrong. Please try again.")
				// Reset credentials to prompt again
				request.Password = ""
				request.Totp = ""
				request.Email = "" // Also reset email to prompt again
				continue
			}
			// For other errors, return immediately
			return fmt.Errorf("failed to login: %v", err)
		}

		// Parse response
		var response LoginResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("failed to parse login response: %v", err)
		}

		clt.accessToken = response.AccessToken
		clt.refreshToken = response.RefreshToken

		return nil
	}
}

func (clt *Client) Refresh(request RefreshTokenRequest) (err error) {
	// Send refresh request
	body, err := clt.doRequest("POST", "/user/refresh", request)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %v", err)
	}
	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse refresh response: %v", err)
	}

	clt.accessToken = response.AccessToken
	clt.refreshToken = response.RefreshToken

	return nil
}

func (clt *Client) Profile(request WithTokenRequest) (err error, userResponse UserResponse) {
	clt.SetAccessToken(request.AccessToken)

	headers := map[string]string{
		"Authorization": clt.GetAccessToken(),
		"Content-Type":  "application/json",
	}

	bodyGetUser, err := clt.doRequestWithHeaders("GET", "/user/profile", nil, headers)
	if err != nil {
		return fmt.Errorf("failed to fetch user profile: %v", err), UserResponse{}
	}

	if err := json.Unmarshal(bodyGetUser, &userResponse); err != nil {
		return fmt.Errorf("failed to parse user profile: %v", err), UserResponse{}
	}

	return nil, userResponse
}

func (clt *Client) UpdateUserPassword(request UpdateUserPasswordRequest) (err error) {
	// Send request
	_, err = clt.doRequest("PATCH", "/user/password", request)
	if err != nil {
		return
	}

	return
}
