'use strict'

/**
  * File {
  *   id: number,
  *   mimetype: string,
  *   url: string,        // relative file url (hash + .ext), without /uploads/
  * }
  *
  * window.files: File[]
  * window.currentFile: number = file.ID | undefined
  */

/**
  * Show file viewer.
  * @param url - relative file url (hash + .ext)
  */
function showFileViewer(fileID) {
  const index = window.files.findIndex(item => item.id === fileID)
  if (index === -1) {
    console.error("Wrong file id, maybe you forgot to set window.files?")
    return
  }

  const file = window.files[index]

  _setFile(file.id, file.url, file.mimetype)
}

function hideImageViewer() {
  const fileViewer = document.getElementById("fileviewer")
  if (!fileViewer) {
    console.error("File viewer container was not found")
    return
  }

  const fileViewerContent = document.getElementById("fileviewer-content")
  if (!fileviewer) {
    console.error("File viewer content container was not found")
    return
  }

  fileViewer.style.display = "none"
  fileViewerContent.innerHTML = ""
}

function showNextFile(e) {
  e.stopPropagation()

  const fileID = window.currentFile
  if (!fileID) return

  const index = window.files.findIndex(item => item.id === fileID)
  if (index === -1) {
    console.error("Wrong file id")
    return
  }

  if (index + 1 < window.files.length) {
    const file = window.files[index + 1]
    _setFile(file.id, file.url, file.mimetype)
  } else {
    const file = window.files[0]
    _setFile(file.id, file.url, file.mimetype)
  }
}

function showPrevFile(e) {
  e.stopPropagation()

  const fileID = window.currentFile
  if (!fileID) return

  const index = window.files.findIndex(item => item.id === fileID)
  if (index === -1) {
    console.error("Wrong file id")
    return
  }

  if (index - 1 >= 0) {
    const file = window.files[index - 1]
    _setFile(file.id, file.url, file.mimetype)
  } else {
    const file = window.files[window.files.length - 1]
    _setFile(file.id, file.url, file.mimetype)
  }
}

function _setFile(fileID, url, mimetype) {
  const fileviewer = document.getElementById("fileviewer")
  if (!fileviewer) {
    console.error("File viewer container was not found")
    return
  }
  const fileviewerContent = document.getElementById("fileviewer-content")
  if (!fileviewer) {
    console.error("File viewer content container was not found")
    return
  }

  fileviewer.style.display = "flex"

  if (mimetype.startsWith("image")) {
    fileviewerContent.innerHTML = `<img src="/uploads/${ url }" />`
    window.currentFile = fileID
  } else if (mimetype.startsWith("video")) {
    fileviewerContent.innerHTML = `<video src="/uploads/${ url }" autoplay controls />`
    window.currentFile = fileID
  } else if (mimetype.startsWith("audio")) {
    fileviewerContent.innerHTML = `<audio src="/uploads/${ url }" autoplay controls />`
    window.currentFile = fileID
  }
}

function _bindEvents() {
  const fileViewer = document.getElementById("fileviewer")
  if (!fileViewer) {
    console.error("File viewer container was not found")
    return
  }

  fileViewer.addEventListener('click', hideImageViewer)
}

window.addEventListener('load', _bindEvents)
