import { createResource } from 'solid-js'
import type { Component } from 'solid-js'

import { fetchBoards } from '../../api'

const BoardsPage: Component = () => {
    const [boards] = createResource(fetchBoards)

    return (
        <div>
        </div>
    )
}

export default BoardsPage
