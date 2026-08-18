# File Storage

Galaplate supports multiple file storage providers through a unified interface: Local, S3, Google Cloud Storage, and Google Drive.

## Configuration

Configure providers in `config/filesystems.yaml`:

```yaml
default: ${FILESYSTEM_DRIVER:local}
max_size: ${FILESYSTEM_MAX_SIZE:10485760}

allowed_types:
  - image/jpeg
  - image/png
  - application/pdf

disks:
  local:
    driver: local
    path: ${FILESYSTEM_LOCAL_PATH:storage/app/uploads}

  s3:
    driver: s3
    region: ${AWS_REGION:}
    bucket: ${AWS_BUCKET:}
    key: ${AWS_ACCESS_KEY_ID:}
    secret: ${AWS_SECRET_ACCESS_KEY:}
    endpoint: ${S3_ENDPOINT:}
    use_path_style_endpoint: ${AWS_USE_PATH_STYLE:false}
    path_prefix: ${AWS_PATH_PREFIX:}

  gcs:
    driver: gcs
    project_id: ${GCP_PROJECT:}
    bucket: ${GCS_BUCKET:}
    key_file: ${GOOGLE_APPLICATION_CREDENTIALS:}

  google_drive:
    driver: google_drive
    service_account_file: ${GOOGLE_SERVICE_ACCOUNT_FILE:}
    folder_id: ${GOOGLE_DRIVE_FOLDER_ID:}
```

## Uploading Files

```go
import filestorage "github.com/galaplate/core/file-storage"
import "github.com/galaplate/core/file-storage/factory"

func (c *Controller) Upload(ctx fiber.Ctx) error {
    file, err := ctx.FormFile("document")
    if err != nil {
        return err
    }

    // Upload using default provider
    metadata := factory.Upload(file)

    // Upload using specific provider
    metadata := factory.UploadWith("s3", file)

    if metadata.Error != "" {
        return fmt.Errorf("upload failed: %s", metadata.Error)
    }

    return ctx.JSON(fiber.Map{
        "filename": metadata.FileName,
        "path":     metadata.FilePath,
        "size":     metadata.FileSize,
    })
}
```

## Download URLs

```go
url := factory.GetDownloadURL(metadata)
```

## Deleting Files

```go
err := factory.Delete(metadata.FilePath, metadata.StorageType)
```

## Checking Existence

```go
exists := factory.ValidateFileExists(metadata)
```

## Switching Providers

```go
// Change default at runtime
factory.SetProvider("s3")
```

## Upload Metadata

```go
type UploadMetadata struct {
    FileName      string  // Sanitized filename
    FilePath      string  // Storage path or identifier
    FileSize      int64   // Size in bytes
    MimeType      string  // MIME type
    StorageType   string  // Provider name
    GoogleDriveID *string // GDrive file ID (if applicable)
    Error         string  // Error message if failed
}
```
