import { h } from "preact";
import { useState } from "preact/hooks";

export interface AppProps {
    name: string;
}

export default function App(props: AppProps) {
    const [count, setCount] = useState(0);
    return (
        <section>
            <header>
                <h2>Hello, {props.name}</h2>
            </header>
            <p>{count}</p>
            <button onClick={() => {setCount(count + 1)}}>Increase</button>
        </section>
    );
}
