import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { ROOM_WEBSOCKET_SERVER_ORIGIN } from "./constants";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<React.StrictMode>
		<App />
	</React.StrictMode>
);
