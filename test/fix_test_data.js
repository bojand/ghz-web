#!/usr/bin/env node

const fs = require('fs')
const path = require('path')

const files = fs.readdirSync(__dirname)

const runFiles = files.map(p => {
  const ext = path.extname(p)
  const bn = path.basename(p)
  if (ext === '.json' && bn.indexOf('run') >= 0) {
    const inputPath = path.join(__dirname, p)
    return inputPath
  }
})

let n = 0
const MONTH = (new Date()).getMonth()
runFiles.forEach(rf => {
  try {
    if (!rf) {
      return
    }

    n++
    if (n > 27) {
        console.log('maximum reached skipping...')
        return
    }


    const content = fs.readFileSync(rf, 'utf8')
    const data = JSON.parse(content)
    const date = new Date()
    date.setMonth(MONTH)
    date.setDate(n)
    data.date = date.toISOString()
    fs.writeFileSync(rf, JSON.stringify(data), 'utf8')
  } catch (e) {
    console.log(e)
  }
})
