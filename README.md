# deviceshifu-ai-sfn

1. Install YoMo cli

```sh
curl -fsSL https://get.yomo.run | sh
```

2. Start mock server

```sh
uvicorn mock_server:app --port 30080
```

3. vivgrid.com

https://dashboard.vivgrid.com

set system prompt

```text
You are a helpful embodied intelligence assistant.Use the supplied tools to assist the user.If the final question is related to the current status of LED and PLC, please look up from the chat history first; and you should call the function "get-image" only if there's no clear answer in chat history.
```

```sh
export VIVGRID_TOKEN="********"
```

4. Run get_image SFN

create a new project with no tools from vivgrid

```sh
cd get_image

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"
export VIVGRID_TOKEN_WITHOUT_TOOLS="******"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "messages": [{"role": "user", "content": "Hi, can you tell what the camera is seeing?"}]
   }'
```

5. Run set_plc_output SFN

```sh
cd set_plc_output

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "messages": [{"role": "user", "content": "Can you set the PLC output to true?"}]
   }'
```

6. Run set_led SFN

```sh
cd set_led

export YOMO_SFN_CREDENTIAL="app-key-secret:$VIVGRID_TOKEN"

yomo run app.go
```

```sh
curl https://openai.vivgrid.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $VIVGRID_TOKEN" \
  -d '{
     "messages": [{"role": "user", "content": "Can you set the display number on the LED to 4005?"}]
   }'
```

7. Run the whole process of the chat demo

```sh
jupyter-notebook chat.ipynb
```
