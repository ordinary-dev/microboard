import { For, Show, createResource } from 'solid-js'
import { useParams, useNavigate } from '@solidjs/router'
import type { Component } from 'solid-js'

import { fetchBoard } from '../../api'
import { createThread, fetchThreads } from '../../api'
import styles from './index.module.css'

const Board: Component = () => {
    const params = useParams()

    return (
        <div class={styles.Page}>
            <BoardTitle boardCode={ params.code } />
            <ListOfThreads boardCode={ params.code } />
            <NewThreadForm boardCode={ params.code } />
        </div>
    )
}

const BoardTitle: Component<{ boardCode: string }> = (params) => {
    const [board] = createResource(params.boardCode, fetchBoard)

    return (
        <div class={styles.Title}>
            <Show when={board() !== undefined}>
                /{ board().code }/ - { board().name }
            </Show>
        </div>
    )
}

const ListOfThreads: Component<{ boardCode: string }> = (params) => {
    const [threads] = createResource(() => fetchThreads(params.boardCode, 10, 0))

    return (
        <div>
            <For each={threads()} fallback={<div>Loading...</div>}>
                {thread => <div>#{thread.id}</div>}
            </For>
        </div>
    )
}

const NewThreadForm: Component<{ boardCode: string }> = (props) => {
    const navigate = useNavigate()
    
    return (
        <form class={styles.NewThreadForm} onSubmit={(e) => submitThread(e, navigate)}>
            <h2>New thread</h2>
            <input
                type="hidden"
                name="boardCode"
                value={props.boardCode}
                readonly
            />
            <textarea
                name="body"
                placeholder="Text"
                required
            />
            <button>Create</button>
        </form>
    )
}

interface NewThreadForm extends HTMLFormElement {
    boardCode: HTMLInputElement
    body: HTMLTextAreaElement
}

async function submitThread(e: Event, navigate: (to: string) => void) {
    e.preventDefault()
    const target = e.target as NewThreadForm

    try {
        const thread = await createThread(target.boardCode.value, target.body.value)
        navigate(`/threads/${thread.id}`)
    }
    catch (err) {
        alert(err.message)
    }
}

export default Board
