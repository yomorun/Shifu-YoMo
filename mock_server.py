from fastapi import FastAPI
from fastapi.responses import FileResponse

app = FastAPI()


@app.get('/deviceshifu-camera/capture')
async def camera():
    return FileResponse('./image.png')


@app.get('/deviceshifu-plc/sendsinglebit')
async def plc():
    return ''


@app.post('/deviceshifu-led/number')
async def led():
    return ''
