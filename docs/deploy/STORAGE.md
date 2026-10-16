# Storage integrations

## S3 / AWS

The backend reads these environment variables:

- `S3_BUCKET` or `S3_BUCKET_NAME`
- `AWS_RESOURCE_ARN` may provide the bucket name if the bucket variable is absent
- `S3_REGION` or `AWS_REGION`
- `S3_ACCESS_KEY_ID` and `S3_SECRET_ACCESS_KEY` for static credentials
- `AWS_ROLE_ARN` / `AWS_ACCOUNT_ID` can be used by external platforms such as Vercel; the Go runtime must still map it to an available credential chain.

Never commit real credentials. Store them in the secret manager or environment.

## Vercel Blob

For frontend or serverless components, use the Vercel Blob client with private access:

```js
import { put } from "@vercel/blob";

const blob = await put("articles/blob.txt", "Hello World!", { access: "private" });
```

`BLOB_READ_WRITE_TOKEN` must be stored in Vercel environment configuration, not in the repository.
