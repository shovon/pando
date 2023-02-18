import React, { Provider } from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { GetRoom } from "./components/GetRoom";
import "./index.css";
import { RoomView } from "./components/Room";

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
