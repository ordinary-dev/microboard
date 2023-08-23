'use strict'

/**
 * Hide the form to create a thread on page load.
 * It is visible by default so that users can create a thread without js.
 */
function hideReplyForm() {
    const replyForm = document.getElementById('replyForm')
    const replyBtn = document.getElementById('newReplyBtn')

    replyBtn.style.display = 'block'
    replyForm.style.display = 'none'
}

function showReplyForm() {
    const replyForm = document.getElementById('replyForm')
    const replyBtn = document.getElementById('newReplyBtn')

    replyBtn.style.display = 'none'
    replyForm.style.display = 'block'
}

window.addEventListener('load', hideReplyForm)
