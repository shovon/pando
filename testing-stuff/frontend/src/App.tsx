import { StrictMode, useState } from "react";
import reactLogo from "./assets/react.svg";
import "./App.css";
import { RoomView } from "./RoomView/RoomView";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { GetRoom } from "./GetRoom";

const router = createBrowserRouter([
	{
		path: "/room/:roomId",
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
