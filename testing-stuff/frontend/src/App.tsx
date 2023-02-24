import { StrictMode, useState } from "react";
import reactLogo from "./assets/react.svg";
import "./App.css";
import { Room } from "./room";
import { RoomView } from "./RoomView";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { GetRoom } from "./GetRoom";

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
		<StrictMode>
			<RouterProvider router={router} />
		</StrictMode>
	);
}

export default App;
