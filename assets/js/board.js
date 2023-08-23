'use strict'

/**
 * Hide the form to create a thread on page load.
 * It is visible by default so that users can create a thread without js.
 */
function hideThreadForm() {
    const threadForm = document.getElementById('threadForm')
    const threadBtn = document.getElementById('newThreadBtn')

    threadBtn.style.display = 'block'
    threadForm.style.display = 'none'
}

function showThreadForm() {
    const threadForm = document.getElementById('threadForm')
    const threadBtn = document.getElementById('newThreadBtn')

    threadBtn.style.display = 'none'
    threadForm.style.display = 'block'
}

window.addEventListener('load', hideThreadForm)
