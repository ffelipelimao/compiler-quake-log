package reader

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockCreate(data string) (*os.File, error) {
	file, err := os.CreateTemp("", "test_file.txt")
	if err != nil {
		return nil, err
	}
	_, err = file.WriteString(data)
	if err != nil {
		file.Close()
		return nil, err
	}
	return file, nil
}

func Test_LoadRows(t *testing.T) {
	tests := []struct {
		name          string
		inputData     string
		expectedRows  []string
		expectedError error
	}{
		{
			name:          "should load lines successfully",
			inputData:     "Line 1\nLine 2\nLine 3\n",
			expectedRows:  []string{"Line 1", "Line 2", "Line 3"},
			expectedError: nil,
		},
		{
			name:          "should return an error when file does not exist",
			inputData:     "",
			expectedRows:  nil,
			expectedError: errors.New("error to open file open nonexistentfile.txt: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile, err := mockCreate(tt.inputData)
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name())
			defer tempFile.Close()

			if tt.expectedError != nil {
				reader := New("nonexistentfile.txt")
				_, err := reader.LoadRows()
				assert.EqualError(t, err, tt.expectedError.Error())

			} else {
				reader := New(tempFile.Name())
				rows, err := reader.LoadRows()
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, tt.expectedRows, rows)
			}

		})
	}
}
