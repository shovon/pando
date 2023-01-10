import React, { Provider } from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { GetRoom } from "./components/GetRoom";
import "./index.css";
import { RoomView } from "./components/Room";
// import { generateKeys } from "@sparkscience/wskeyid-browser/src/utils";
// import { Session } from "./session";
// import { toAsyncIterable } from "@sparkscience/wskeyid-browser/src/pub-sub";

const router = createBrowserRouter([
	{
		path: "/room/:id",
		element: <RoomView />,
	},
	{
		path: "*",
		element: <GetRoom />,
	},
]);

// Promise.resolve()
// 	.then(async function () {
// 		const keys = await generateKeys();
// 		const session = await new Session(
// 			"ws://localhost:8080/room/some_room",
// 			keys
// 		);

// 		session.sessionStatusChangeEvents.addEventListener((status) => {
// 			if (status.type === "CONNECTING") {
// 				console.log(
// 					"status type is %s and sub status is %s",
// 					status.type,
// 					status.status
// 				);
// 			} else {
// 				console.log("Status type is %s", status.type);
// 			}
// 		});

// 		for await (const event of toAsyncIterable(session.messageEvents)) {
// 			console.log(event);
// 		}
// 	})
// 	.catch(console.error);

function App() {
	return (
		<React.StrictMode>
			<RouterProvider router={router} />
		</React.StrictMode>
	);
}

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<App />
);
