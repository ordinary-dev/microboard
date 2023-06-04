import { For, createResource } from 'solid-js'
import type { Component } from 'solid-js'

import { fetchBoards, createBoard, deleteBoard } from '../../api'
import type { Board } from '../../api'
import styles from './index.module.css'

const AdminPanel: Component = () => {
    const [boards, { refetch }] = createResource(fetchBoards)
    return (
        <div class={styles.AdminPanel}>
            <h1>Boards</h1>
            <form class={styles.BoardForm} onSubmit={(e) => submitNewBoard(e).then(refetch)}>
                <input
                    placeholder="Code"
                    name="boardCode"
                    required
                />
                <input
                    placeholder="Name"
                    name="boardName"
                    required
                />
                <button>Create</button>
            </form>
            <For each={boards()} fallback={<div>Loading...</div>}>
                {board => <BoardRecord
                    code={board.code}
                    name={board.name}
                    onDelete={() => deleteBoard(board.code).then(refetch)}
                />}
            </For>
        </div>
    )
}

interface BoardForm extends HTMLFormElement {
    boardCode: HTMLInputElement
    boardName: HTMLInputElement
}

async function submitNewBoard(e: Event) {
    e.preventDefault()
    const target = e.target as BoardForm
    const { boardCode, boardName } = target
    try {
        await createBoard(boardCode.value, boardName.value)
        target.reset()
    }
    catch(err) {
        console.error(err)
        alert('Failed to create a board, see console for details')
    }
}

interface BoardRecordProps extends Board {
    onDelete: () => void,
}

const BoardRecord: Component<BoardRecordProps> = (props) => {
    return (
        <div class={styles.BoardRecord}>
            <div class={styles.BoardCode}>{ props.code }</div>
            <div class={styles.BoardName}>{ props.name }</div>
            <button onClick={props.onDelete} class={styles.DeleteBoardBtn}>x</button>
        </div>
    )
}

export default AdminPanel
