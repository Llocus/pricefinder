import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
//import reportWebVitals from './reportWebVitals';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import Home from './screens/home';
import Products from './screens/products';
import { defaultSystem } from "@chakra-ui/react"
import { ChakraProvider } from '@chakra-ui/react';
import { IconContext } from "react-icons";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
    //element: <Products />,
  },
  {
    path: "/search/:search",
    element: <Products />,
  },
]);

root.render(
  <IconContext.Provider value={{ className: "global-class-name" }}>
    <ChakraProvider value={defaultSystem}>
      <RouterProvider router={router} />
    </ChakraProvider>
  </IconContext.Provider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
//reportWebVitals();
