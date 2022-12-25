import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import App from "./App";
import { GetRoom } from "./GetRoom";
import "./index.css";
import { Room } from "./Room";

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
