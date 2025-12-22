# Deployment Guide - Render.com

This guide walks you through deploying the Image Analyzer application on Render.com with your custom domain.

## Prerequisites

- GitHub account with both repositories
- Render.com account (free, no credit card required initially)
- Google Cloud credentials JSON file
- Domain pedro00627.com configured in Hostinger

## Step 1: Prepare Google Cloud Credentials

1. Get your credentials JSON content ready (from `C:\gcp\credentials.json`)
2. You'll need to add this as an environment variable in Render

## Step 2: Deploy Backend API

### 2.1 Create Web Service

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click **"New +"** → **"Web Service"** (NOT "Blueprint" or "Static Site")
3. Connect your GitHub account if not already connected
4. Select repository: `pedro00627/image-analyzer-api`
5. Select branch: `develop` (or `main` after merging)

> **Important:** Create as **Web Service**, not Blueprint. The render.yaml file is for reference only.

### 2.2 Configure Service

**Basic Settings:**
- **Name:** `image-analyzer-api`
- **Region:** Oregon (or closest to your users)
- **Branch:** `develop`
- **Runtime:** Docker
- **Instance Type:** Free

**Environment Variables:**

Click "Add Environment Variable" for each:

```
AI_PROVIDER=google_vision
PORT=8080
ALLOWED_ORIGINS=https://pero00627.com,https://www.pero00627.com
MAX_FILE_SIZE=10485760
ALLOWED_FILE_TYPES=image/jpeg,image/png,image/gif,image/webp
AI_SERVICE_TIMEOUT=30s
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=10s
HTTP_IDLE_TIMEOUT=60s
```

**GOOGLE_APPLICATION_CREDENTIALS (Important!):**

This requires TWO configurations:

**Step A: Add Secret File**
1. Click **"Add Secret File"** (not regular environment variable)
2. **Filename:** `credentials.json` (just the name, no path)
3. **Contents:** Open your `C:\gcp\credentials.json` file and copy-paste the ENTIRE JSON content

Example of what to paste:
```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...",
  ...
}
```

**Step B: Add Environment Variable**
4. Click **"Add Environment Variable"**
5. **Key:** `GOOGLE_APPLICATION_CREDENTIALS`
6. **Value:** `/etc/secrets/credentials.json`

> **How it works:** Render creates the file at `/etc/secrets/credentials.json` and your environment variable tells the app where to find it.

### 2.3 Deploy

1. Click **"Create Web Service"**
2. Wait for build and deployment (5-10 minutes)
3. Once deployed, note your URL: `https://image-analyzer-api.onrender.com`
4. Test health endpoint: `https://image-analyzer-api.onrender.com/health`

## Step 3: Deploy Frontend

### 3.1 Create Static Site

1. From Render Dashboard, click **"New +"** → **"Static Site"** (NOT "Web Service")
2. Select repository: `pedro00627/image-analyzer-web`
3. Select branch: `develop`

> **Important:** Create as **Static Site** for the frontend, not Web Service.

### 3.2 Configure Static Site

**Basic Settings:**
- **Name:** `image-analyzer-web`
- **Branch:** `develop`
- **Build Command:** `npm install && npm run build`
- **Publish Directory:** `dist/image-analyzer-web`

**Environment Variables:**
```
NODE_VERSION=18
```

### 3.3 Deploy

1. Click **"Create Static Site"**
2. Wait for build (5-8 minutes)
3. Once deployed, your site will be at: `https://image-analyzer-web.onrender.com`

## Step 4: Configure Custom Domains

### 4.1 Configure Backend Domain (api.pedro00627.com)

1. Go to your `image-analyzer-api` web service in Render
2. Click **"Settings"** → **"Custom Domains"**
3. Click **"Add Custom Domain"**
4. Enter: `api.pedro00627.com`
5. Render will show DNS records to configure

### 4.2 Configure Frontend Domain (pedro00627.com)

1. Go to your `image-analyzer-web` static site
2. Click **"Settings"** → **"Custom Domains"**
3. Click **"Add Custom Domain"**
4. Add both:
   - `pedro00627.com`
   - `www.pedro00627.com`
5. Render will show DNS records for each

### 4.3 Configure DNS in Hostinger

1. Log into Hostinger
2. Go to **Domains** → Select `pedro00627.com`
3. Go to **DNS / Name Servers**
4. Add these records:

**For API subdomain (Backend):**
```
Type: CNAME
Name: api
Value: image-analyzer-api.onrender.com
TTL: 3600
```

**For root domain (Frontend):**
```
Type: A
Name: @
Value: [IP address provided by Render]
TTL: 3600
```

**For www subdomain (Frontend):**
```
Type: CNAME
Name: www
Value: image-analyzer-web.onrender.com
TTL: 3600
```

5. Save changes
6. Wait for DNS propagation (5-60 minutes)

### 4.4 Verify Domains

After DNS propagation, test all URLs:
- ✅ Frontend: `https://pedro00627.com`
- ✅ Frontend: `https://www.pedro00627.com`
- ✅ Backend: `https://api.pedro00627.com/health`

## Step 5: Update Backend CORS

Once your domains are configured, update backend CORS to allow your custom domains:

1. Go to backend service in Render Dashboard
2. Go to **"Environment"** tab
3. Update `ALLOWED_ORIGINS` environment variable to:
```
ALLOWED_ORIGINS=https://pedro00627.com,https://www.pedro00627.com
```
4. Click **"Save Changes"**
5. Service will automatically redeploy (takes ~2-3 minutes)

> **Note:** Remove localhost origins in production for better security.

## Step 6: Enable SSL (Automatic)

Render automatically provides free SSL certificates via Let's Encrypt. Once DNS is configured:

1. Render detects custom domain
2. Automatically provisions SSL certificate
3. Forces HTTPS redirects

## Testing Deployment

### Test Backend
```bash
# Health check
curl https://image-analyzer-api.onrender.com/health

# Test analyze endpoint (with image file)
curl -X POST https://image-analyzer-api.onrender.com/api/analyze \
  -F "image=@/path/to/test-image.jpg"
```

### Test Frontend
1. Open `https://pero00627.com`
2. Upload an image
3. Verify analysis works

## Important Notes

### Free Tier Limitations

**Backend (Web Service - Free):**
- ✅ 750 hours/month (enough for 24/7 operation)
- ⚠️ Spins down after 15 minutes of inactivity
- ⚠️ First request after spin-down takes ~30 seconds (cold start)

**Frontend (Static Site - Free):**
- ✅ Unlimited bandwidth
- ✅ Always fast, no cold starts
- ✅ CDN included

### Cold Start Solution

To keep backend warm (optional):
- Use a free uptime monitor (UptimeRobot, Freshping)
- Ping `/health` every 10 minutes
- Prevents spin-down

### Monitoring

1. Go to service dashboard
2. View logs in real-time
3. Check metrics (requests, response times)

## Troubleshooting

### Backend won't start
- Check logs in Render dashboard
- Verify Google credentials are properly configured
- Ensure all environment variables are set

### Frontend shows errors
- Check browser console
- Verify API URL in environment.prod.ts matches backend URL
- Check CORS settings in backend

### Domain not resolving
- Verify DNS records in Hostinger
- Wait up to 48 hours for full propagation
- Use DNS checker: https://dnschecker.org/

### API requests failing
- Check CORS settings include your domain
- Verify backend is running (check /health endpoint)
- Check browser network tab for specific errors

## Costs

**Current setup: $0/month**

**If you need to upgrade later:**
- Backend always-on (no cold starts): $7/month
- Custom domain SSL on backend: Included free

## Next Steps

After successful deployment:

1. ✅ Set up uptime monitoring
2. ✅ Configure error tracking (optional: Sentry)
3. ✅ Set up GitHub Actions for auto-deployment
4. ✅ Monitor usage and performance

## Useful Links

- [Render Dashboard](https://dashboard.render.com/)
- [Render Docs - Docker](https://render.com/docs/docker)
- [Render Docs - Static Sites](https://render.com/docs/static-sites)
- [Render Docs - Custom Domains](https://render.com/docs/custom-domains)
