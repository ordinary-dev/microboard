export interface Thread {
    id: number
    body: string
    thread_id: number
    created_at: string
}

export async function fetchThreads(boardCode: string, limit: number, offset: number): Promise<Thread[]> {
    const res = await fetch(`/api/v0/threads/by_board_code/${boardCode}?limit=${limit}&offset=${offset}`)
    if (!res.ok)
        throw new Error('Failed to fetch the list of threads', { cause: res })
    return res.json()
}

export async function createThread(boardCode: string, body: string): Promise<Thread> {
    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            board_code: boardCode,
            body,
        }),
    }
    const res = await fetch('/api/v0/threads', options)
    if (!res.ok)
        throw new Error('Failed to create a new thread', { cause: res })
    return res.json()
}
