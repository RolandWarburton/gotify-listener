# gotify-listen

Forwards gotify messages via `notify-send`.

Ensure you have `/etc/gotify/cli.json` configured.

```json
{
  "url": "https://push.example.net",
  "token": "D2ESUSbot.8iHDz"
}
```

## Testing

Generate an APP token in the gotify web UI to authorize pushing notifications.

```bash
GOTIFY_TOKEN=a8ms-A8ADanQR7G1 gotify push -t "Test" "Hello world"
```
