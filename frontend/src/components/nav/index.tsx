import { Component } from 'solid-js'
import { A } from '@solidjs/router'

import styles from './index.module.css'


export const NavBar: Component = () => {
    return (
        <div class={styles.NavBar}>
            <A href='/' activeClass={styles.ActiveLink} end>
                Home
            </A>
            <A href='/admin' activeClass={styles.ActiveLink}>
                Admin panel
            </A>
        </div>
    )
}
