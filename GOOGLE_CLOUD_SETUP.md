# Google Cloud Authentication with Docker

There are several ways to authenticate with Google Cloud services from Docker containers:

## Option 1: Service Account Key File (Recommended for Development)

1. **Create a Service Account Key**:
   ```bash
   # Download the key file from Google Cloud Console
   # Or use gcloud CLI:
   gcloud iam service-accounts keys create gcp-key.json \
     --iam-account=your-service-account@your-project.iam.gserviceaccount.com
   ```

2. **Add to `.env` file**:
   ```env
   GOOGLE_APPLICATION_CREDENTIALS=./gcp-key.json
   ```

3. **Add to `.gitignore`** (already included):
   ```
   gcp-key.json
   *.json
   ```

4. **Docker Compose** will automatically mount the file (already configured in `docker-compose.yml`)

## Option 2: Application Default Credentials (ADC)

If you're running Docker locally and have `gcloud` CLI installed:

1. **Authenticate locally**:
   ```bash
   gcloud auth application-default login
   ```

2. **Mount the credentials directory**:
   ```yaml
   volumes:
     - ~/.config/gcloud:/root/.config/gcloud:ro
   ```

3. **Set environment variable**:
   ```env
   GOOGLE_APPLICATION_CREDENTIALS=/root/.config/gcloud/application_default_credentials.json
   ```

## Option 3: Environment Variable (Service Account JSON)

1. **Export the key as environment variable**:
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS_JSON=$(cat gcp-key.json)
   ```

2. **Update docker-compose.yml** to use it:
   ```yaml
   environment:
     GOOGLE_APPLICATION_CREDENTIALS_JSON: ${GOOGLE_APPLICATION_CREDENTIALS_JSON}
   ```

3. **In your Go code**, you can read from the env var:
   ```go
   import "google.golang.org/api/option"
   
   // In your code:
   ctx := context.Background()
   var opts []option.ClientOption
   if json := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"); json != "" {
       opts = append(opts, option.WithCredentialsJSON([]byte(json)))
   }
   client, err := speech.NewClient(ctx, opts...)
   ```

## Option 4: Workload Identity (Production - GKE/GCP)

For production on Google Cloud (GKE, Cloud Run, etc.), use Workload Identity:

1. **No credentials needed** - GCP automatically injects credentials
2. **Service account** is attached to the workload
3. **No code changes needed** - just use default credentials

## Recommended Setup for Development

1. Create a service account key file:
   ```bash
   gcp-key.json
   ```

2. Add to `.env`:
   ```env
   GOOGLE_APPLICATION_CREDENTIALS=./gcp-key.json
   ```

3. The docker-compose.yml already mounts this file automatically

## Security Notes

- **Never commit** `gcp-key.json` to git
- Use **least privilege** - only grant necessary permissions
- For production, use **Workload Identity** instead of key files
- Rotate keys regularly

