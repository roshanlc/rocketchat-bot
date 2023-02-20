// parsing command line data
const { argv } = require('node:process')

// argv.includes("-d") ? j = console.info(argv[argv.indexOf("-d") + 1]) : console.info("Does not contain '-d' switch")
let data = {}
if (argv.includes("-d")) {
    data = JSON.parse(argv[argv.indexOf("-d") + 1])
}
let replyData = {
    "text_reply": "Does not contain any data in the arguement",
    "image_url": "",
    "video_url": "",
}

if (data != {}) {
    replyData = {
        "text_reply": `Msg recieved =${data.msg}`,
        "image_url": "",
        "video_url": "",
    }
}

console.info(JSON.stringify(replyData))