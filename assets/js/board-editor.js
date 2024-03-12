'use strict'

async function deleteBoard(boardCode) {
  if (!confirm(`Are you sure you want to delete /${boardCode}/?`)) return

  try {
    const options = {
      method: "DELETE",
    }
    const res = await fetch(`/api/v0/boards/${boardCode}`, options)

    if (!res.ok) {
      alert(`Something went wrong, code: ${res.status}.`)
      return
    }

    location.reload()
  }
  catch(err) {
    alert(err)
  }
}
