package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

const tokenFile = "token.json"

type GDriveClient struct {
	Service *drive.Service
	Config *oauth2.Config
}


func NewGDriveClient(credentialsFile string) (*GDriveClient, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	// Get token and client
	client := getClient(config)

	// Create Drive service
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Drive service: %v", err)
	}

	return &GDriveClient{Service: srv, Config: config}, nil
}

func (c *GDriveClient) UploadFile(filePath string, mimeType string, parentFolderID string) (*drive.File, error) {
	file := &drive.File{
		Title:    filePath,
		MimeType: mimeType,
		Parents:  []*drive.ParentReference{{Id: parentFolderID}},
	}
	createdFile, err := c.Service.Files.Insert(file).Media(nil).Do()
	if err != nil {
		return nil, err
	}
	return createdFile, nil
}

func (c *GDriveClient) CreateFolder(folderName string, parentFolderID string) (*drive.File, error) {
	folder := &drive.File{
		Title:    folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []*drive.ParentReference{{Id: parentFolderID}},
	}
	createdFolder, err := c.Service.Files.Insert(folder).Do()
	if err != nil {
		return nil, err
	}
	return createdFolder, nil
}

func (c *GDriveClient) ListFilesInFolder(folderID string) ([]*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	fileList, err := c.Service.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}
	return fileList.Items, nil
}


func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Open this URL in your browser:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	fmt.Scan(&authCode)

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		panic(fmt.Errorf("failed to retrieve token: %v", err))
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving token to %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to save token: %v\n", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}