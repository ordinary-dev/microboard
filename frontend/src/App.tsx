import { lazy } from 'solid-js'
import type { Component } from 'solid-js'
import { Routes, Route } from '@solidjs/router'

import { NavBar } from './components/nav'
import styles from './App.module.css'

const BoardListPage = lazy(() => import('./pages/boardList'))
const AdminPanel = lazy(() => import('./pages/admin'))
const BoardPage = lazy(() => import('./pages/board'))

const App: Component = () => {
    return (
        <div>
            <NavBar />
            <Routes>
                <Route path="/" component={ BoardListPage } />
                <Route path="/admin" component={ AdminPanel } />
                <Route path="/boards/:code" component={ BoardPage } />
            </Routes>
        </div>
    )
}

export default App;
