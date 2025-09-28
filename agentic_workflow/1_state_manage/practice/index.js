const express = require("express")
const app = express()

const port = 7000
const fsm = require('./p_1');



app.listen(port, () => {
    console.log('server run at ', port)
})


