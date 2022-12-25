import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { GetRoom } from "./components/GetRoom";
import "./index.css";
import { Room } from "./components/Room";

const router = createBrowserRouter([
	{
		path: "/room/:id",
		element: <Room />,
	},
	{
		path: "*",
		element: <GetRoom />,
	},
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>
);
