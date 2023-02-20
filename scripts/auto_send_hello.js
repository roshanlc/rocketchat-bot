const d = new Date()
const replyData = {
    "text_reply": `Hello from an automated script at unix time: ${d.getTime()}`,
    "image_url": "",
    "video_url": "",
}

console.info(JSON.stringify(replyData))