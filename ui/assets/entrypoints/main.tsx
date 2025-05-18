import { h, render } from "preact";
import App, {type AppProps} from "../App";

const root = document.getElementById("app");
if (!root) {
    throw new Error("app root not found");
}

const props: AppProps = JSON.parse(decodeURIComponent(root.getAttribute("data-props")!));

render(<App name={props.name} />, document.getElementById("app")!)
