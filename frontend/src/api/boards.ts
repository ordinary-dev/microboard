export interface Board {
    code: string
    name: string
}

export async function fetchBoards(): Promise<Board[]> {
    const res = await fetch('/api/v0/boards')
    if (!res.ok)
        throw new Error('Failed to download the list of boards', { cause: res })
    return res.json()
}

export async function createBoard(code: string, name: string): Promise<Board> {
    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            code,
            name,
        }),
    }
    const res = await fetch('/api/v0/boards', options)
    if (!res.ok)
        throw new Error('Failed to create a board', { cause: res })
    return res.json()
}

export async function deleteBoard(code: string): Promise<void> {
    const options = {
        method: 'DELETE',
    }
    const res = await fetch(`/api/v0/boards/${code}`, options)
    if (!res.ok)
        throw new Error('Failed to delete the board', { cause: res })
}
