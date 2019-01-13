import React from 'react';
import Auth from "../Auth";
import Stats from "../Stats"
import {BrowserRouter as Router, Route} from "react-router-dom";

class App extends React.Component {
    render() {
        return <Router>
            <div>
                <Route exact path="/" component={Auth} />
                <Route path="/stats" component={Stats} />
            </div>
        </Router>
    }
}

export default App;
