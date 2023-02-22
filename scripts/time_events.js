// currrent time date object
const date = new Date()

const hour = date.getHours()
const mins = date.getMinutes()

let replyData = {
    "text_reply": `Hello from an automated script at time: ${hour,mins}`,
    "image_url": "",
    "video_url": "",
}

if (hour == 10 && (mins > 10 && mins < 15)){
    replyData.text_reply = "Time for standup .... ihihihihi"
    // send output
    console.info(JSON.stringify(replyData))
    return
}else if (hour == 13 && (mins > 1 && mins < 6)){
    replyData.text_reply = "Khaja khana jaaam .... ihihihi"
    replyData.image_url = "https://raw.githubusercontent.com/roshanlc/rocketchat-bot/haribdr-memes/haribdr-memes/hungry.gif"
    // send output
    console.info(JSON.stringify(replyData))
    return
}else if (hour == 17 && (mins > 5 && mins < 10)){
    replyData.text_reply = "Kati kaam garya, eso ghar ni janu ni .... ihihihi"
    // send output
    console.info(JSON.stringify(replyData))
    return
}
// for empty response send this
replyData = {
    "text_reply": null,
    "image_url": null,
    "video_url": null,
}

console.info(JSON.stringify(replyData))
