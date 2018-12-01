import React from 'react';
import Timer from '../Timer'

const axios = require('axios');

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            app: null
        };
    }

    componentDidMount() {
        axios.get('/api/isauthenticated')
            .then(response => {
                const auth = response.data;
                this.setState({
                    app: <Timer auth={auth}/>
                });
            });
    }

    render() {
        return this.state.app;
    }
}

export default App;
