import { createBrowserRouter } from "react-router-dom";

import App from "../App";
import Home from "../pages/Home";
import Login from "../pages/Login";
import ConfirmSignup from "../pages/ConfirmSignup";
import NewPass from "../pages/NewPass";
import ResetPass from "../pages/ResetPass";
import Pricing from "../pages/Pricing";
import Products from "../pages/Products";
import Admin from "../pages/Admin";
import Checkout from "../pages/Checkout";
import AddBalance from "../pages/AddBalance";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <Error />,
    children: [
      {
        index: true,
        element: <Pricing />,
      },
      {
        path: "/home",
        element: <Home />,
      },
      {
        path: "login",
        element: <Login />,
      },
      {
        path: "confirm",
        element: <ConfirmSignup />,
      },
      {
        path: "newpass",
        element: <NewPass />,
      },
      {
        path: "resetpass",
        element: <ResetPass />,
      },
      {
        path: "products",
        element: <Products />,
      },
      {
        path: "admin",
        element: <Admin />,
      },
      {
        path: "checkout",
        element: <Checkout />,
      },
      {
        path: "balance",
        element: <AddBalance />,
      },
    ],
  },
]);
