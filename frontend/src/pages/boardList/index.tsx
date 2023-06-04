import { For, createResource } from 'solid-js'
import type { Component } from 'solid-js'
import { A } from '@solidjs/router'

import { Board, fetchBoards } from '../../api'
import styles from './index.module.css'

const BoardsPage: Component = () => {
    const [boards] = createResource(fetchBoards)

    return (
        <div class={styles.Page}>
            <Logo />
            <For each={boards()} fallback={<div>Loading...</div>}>
                {board => <BoardRecord
                    code={board.code}
                    name={board.name}
                />}
            </For>
        </div>
    )
}

const Logo: Component = () => {
    return (
        <div class={styles.Logo}>
            Microboard
        </div>
    )
}

const BoardRecord: Component<Board> = (props) => {
    return (
        <A href={`/boards/${props.code}`} class={styles.BoardRecord}>
            <div class={styles.Code}>{ props.code }</div>
            <div class={styles.Name}>{ props.name }</div>
        </A>
    )
}

export default BoardsPage
