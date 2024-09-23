package server_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type GetAccountResponseResult struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Permissions int    `json:"permissions"`
}

type GetAccountResponse struct {
	Result GetAccountResponseResult `json:"result"`
	Error  string                   `json:"error,omitempty"`
}

type CreateAccountParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAccountResponseResult struct {
	ID int `json:"id"`
}

type CreateAccountResponse struct {
	Result CreateAccountResponseResult `json:"result"`
	Error  string                      `json:"error,omitempty"`
}

type ChangeAccountPasswordParams struct {
	Password string `json:"password"`
}

type ChangeAccountPasswordResponseResult struct {
	ID int `json:"id"`
}

type ChangeAccountPasswordResponse struct {
	Result ChangeAccountPasswordResponseResult `json:"result"`
	Error  string                              `json:"error,omitempty"`
}

type DeleteAccountResponseResult struct {
	ID int `json:"id"`
}

type DeleteAccountResponse struct {
	Result DeleteAccountResponseResult `json:"result"`
	Error  string                      `json:"error,omitempty"`
}

func getAccount(url string, client *http.Client, adminToken string, id int) (int, *GetAccountResponse, error) {
	req, err := http.NewRequest("GET", url+"/api/v1/accounts/"+strconv.Itoa(id), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var accountResponse GetAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&accountResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &accountResponse, nil
}

func createAccount(url string, client *http.Client, adminToken string, data *CreateAccountParams) (int, *CreateAccountResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/accounts", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var createResponse CreateAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&createResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &createResponse, nil
}

func changeAccountPassword(url string, client *http.Client, adminToken string, id int, data *ChangeAccountPasswordParams) (int, *ChangeAccountPasswordResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url+"/api/v1/accounts/"+strconv.Itoa(id)+"/change_password", strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var changeResponse ChangeAccountPasswordResponse
	if err := json.NewDecoder(res.Body).Decode(&changeResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &changeResponse, nil
}

func deleteAccount(url string, client *http.Client, adminToken string, id int) (int, *DeleteAccountResponse, error) {
	req, err := http.NewRequest("DELETE", url+"/api/v1/accounts/"+strconv.Itoa(id), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	var deleteResponse DeleteAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&deleteResponse); err != nil {
		return 0, nil, err
	}
	return res.StatusCode, &deleteResponse, nil
}

func TestAccounts(t *testing.T) {
	ts, _, err := setupServer()
	if err != nil {
		t.Fatalf("couldn't create test server: %s", err)
	}
	defer ts.Close()
	client := ts.Client()
	var adminToken string
	var nonAdminToken string
	t.Run("prepare accounts and tokens", prepareAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	t.Run("Get admin account - admin token", func(t *testing.T) {
		statusCode, response, err := getAccount(ts.URL, client, adminToken, 1)
		if err != nil {
			t.Fatalf("couldn't get account: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}
		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}
		if response.Result.ID != 1 {
			t.Fatalf("expected ID 1, got %d", response.Result.ID)
		}
		if response.Result.Username != "testadmin" {
			t.Fatalf("expected username testadmin, got %s", response.Result.Username)
		}
		if response.Result.Permissions != 1 {
			t.Fatalf("expected permissions 1, got %d", response.Result.Permissions)
		}
	})

	t.Run("Get admin account - non admin token", func(t *testing.T) {
		statusCode, response, err := getAccount(ts.URL, client, nonAdminToken, 1)
		if err != nil {
			t.Fatalf("couldn't get account: %s", err)
		}
		if statusCode != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, statusCode)
		}
		if response.Error != "forbidden: admin access required" {
			t.Fatalf("expected error %q, got %q", "forbidden: admin access required", response.Error)
		}
	})

	t.Run("Create account - no password", func(t *testing.T) {
		createAccountParams := &CreateAccountParams{
			Username: "nopass",
			Password: "",
		}
		statusCode, response, err := createAccount(ts.URL, client, adminToken, createAccountParams)
		if err != nil {
			t.Fatalf("couldn't create account: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}
		if response.Result.ID != 3 {
			t.Fatalf("expected ID 3, got %d", response.Result.ID)
		}
	})

	t.Run("Get account", func(t *testing.T) {
		statusCode, response, err := getAccount(ts.URL, client, adminToken, 3)
		if err != nil {
			t.Fatalf("couldn't get account: %s", err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
		}
		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}
		if response.Result.ID != 3 {
			t.Fatalf("expected ID 3, got %d", response.Result.ID)
		}
		if response.Result.Username != "nopass" {
			t.Fatalf("expected username nopass, got %s", response.Result.Username)
		}
		if response.Result.Permissions != 0 {
			t.Fatalf("expected permissions 0, got %d", response.Result.Permissions)
		}
	})

	t.Run("Get account - id not found", func(t *testing.T) {
		statusCode, response, err := getAccount(ts.URL, client, adminToken, 100)
		if err != nil {
			t.Fatalf("couldn't get account: %s", err)
		}
		if statusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, statusCode)
		}
		if response.Error != "Not Found" {
			t.Fatalf("expected error %q, got %q", "Not Found", response.Error)
		}
	})

	t.Run("Create account - no username", func(t *testing.T) {
		createAccountParams := &CreateAccountParams{
			Username: "",
			Password: "password",
		}
		statusCode, response, err := createAccount(ts.URL, client, adminToken, createAccountParams)
		if err != nil {
			t.Fatalf("couldn't create account: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if response.Error != "Username is required" {
			t.Fatalf("expected error %q, got %q", "Username is required", response.Error)
		}
	})

	t.Run("Change account password - success", func(t *testing.T) {
		changeAccountPasswordParams := &ChangeAccountPasswordParams{
			Password: "newPassword1",
		}
		statusCode, response, err := changeAccountPassword(ts.URL, client, adminToken, 1, changeAccountPasswordParams)
		if err != nil {
			t.Fatalf("couldn't create account: %s", err)
		}
		if statusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, statusCode)
		}
		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}
		if response.Result.ID != 1 {
			t.Fatalf("expected ID 1, got %d", response.Result.ID)
		}
	})

	t.Run("Change account password - no user", func(t *testing.T) {
		changeAccountPasswordParams := &ChangeAccountPasswordParams{
			Password: "newPassword1",
		}
		statusCode, response, err := changeAccountPassword(ts.URL, client, adminToken, 100, changeAccountPasswordParams)
		if err != nil {
			t.Fatalf("couldn't create account: %s", err)
		}
		if statusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, statusCode)
		}
		if response.Error != "Not Found" {
			t.Fatalf("expected error %q, got %q", "Not Found", response.Error)
		}
	})

	t.Run("Change account password - no password", func(t *testing.T) {
		changeAccountPasswordParams := &ChangeAccountPasswordParams{
			Password: "",
		}
		statusCode, response, err := changeAccountPassword(ts.URL, client, adminToken, 1, changeAccountPasswordParams)
		if err != nil {
			t.Fatalf("couldn't create account: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if response.Error != "Password is required" {
			t.Fatalf("expected error %q, got %q", "Password is required", response.Error)
		}
	})

	t.Run("Change account password - bad password", func(t *testing.T) {
		changeAccountPasswordParams := &ChangeAccountPasswordParams{
			Password: "password",
		}
		statusCode, response, err := changeAccountPassword(ts.URL, client, adminToken, 1, changeAccountPasswordParams)
		if err != nil {
			t.Fatalf("couldn't create account: %s", err)
		}
		if statusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, statusCode)
		}
		if response.Error != "Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol." {
			t.Fatalf("expected error %q, got %q", "Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol.", response.Error)
		}
	})

	t.Run("Delete account - success", func(t *testing.T) {
		statusCode, response, err := deleteAccount(ts.URL, client, adminToken, 2)
		if err != nil {
			t.Fatalf("couldn't delete account: %s", err)
		}
		if statusCode != http.StatusAccepted {
			t.Fatalf("expected status %d, got %d", http.StatusAccepted, statusCode)
		}
		if response.Error != "" {
			t.Fatalf("expected error %q, got %q", "", response.Error)
		}
		if response.Result.ID != 2 {
			t.Fatalf("expected ID 2, got %d", response.Result.ID)
		}
	})

	t.Run("Delete account - no user", func(t *testing.T) {
		statusCode, response, err := deleteAccount(ts.URL, client, adminToken, 100)
		if err != nil {
			t.Fatalf("couldn't delete account: %s", err)
		}
		if statusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, statusCode)
		}
		if response.Error != "Not Found" {
			t.Fatalf("expected error %q, got %q", "Not Found", response.Error)
		}
	})
}
