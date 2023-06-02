import { lazy } from 'solid-js'
import type { Component } from 'solid-js'
import { Routes, Route } from '@solidjs/router'

import { NavBar } from './components/nav'
import styles from './App.module.css'

const BoardsPage = lazy(() => import('./pages/boards'))
const AdminPanel = lazy(() => import('./pages/admin'))

const App: Component = () => {
    return (
        <div>
            <NavBar />
            <Routes>
                <Route path="/" component={ BoardsPage } />
                <Route path="/admin" component={ AdminPanel } />
            </Routes>
        </div>
    )
}

export default App;
