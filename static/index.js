function format(centis) {
    const hours = Math.floor(centis / (100 * 60 * 60)).toString();
    const minutes = Math.floor(centis % (100 * 60 * 60) / (100 * 60)).toString();
    const seconds = (centis % (100 * 60) / 100).toFixed(2);

    if (hours !== '0') {
        return hours + ':' +
            (minutes.length === 1 ? '0' : '') +
            minutes + ':' +
            (seconds.length === 4 ? '0' : '') +
            seconds;
    } else if (minutes !== '0') {
        return minutes + ':' +
            (seconds.length === 4 ? '0' : '') +
            seconds;
    } else {
        return seconds;
    }
}

class Time extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            running: false,
            timeState: "usual",
            centis: 0
        };

        this.keydown = this.keydown.bind(this);
        this.keyup = this.keyup.bind(this);
        this.start = this.start.bind(this);
        this.stop = this.stop.bind(this);
        document.onkeydown = (e) => {
            if (e.code === "Space") this.keydown(e)
        };
        document.onkeyup = (e) => {
            if (e.code === "Space") this.keyup(e)
        };
    }

    start() {
        this.setState({
            running: true,
            startTime: Date.now()
        });

        this.timerId = setInterval(
            () => this.tick(), 9);
        document.ontouchstart = this.keydown;
    }

    stop() {
        const res = Math.floor((Date.now() - this.state.startTime) / 10);
        clearInterval(this.timerId);
        this.props.saveResult({
            centis: res,
            scramble: this.props.scramble,
            penalty: false,
            datetime: currentDatetime()
        });
        this.setState({
            centis: res,
            running: false
        });
        document.ontouchstart = () => {
        };
    }

    keydown(e) {
        e.preventDefault();
        if (this.state.running) {
            this.stop();
        } else if (isNaN(this.timeoutId)) {
            this.setState({
                timeState: "notReady"
            });
            this.timeoutId = setTimeout(() => {
                this.setState({
                    timeState: "ready"
                });
            }, 300);
        }
    }

    keyup(e) {
        if (!this.state.running) {
            if (this.state.timeState === "ready") {
                this.start()
            }
            clearTimeout(this.timeoutId);
            this.timeoutId = NaN;

            this.setState({
                timeState: "usual"
            });
        }
    }

    tick() {
        this.setState({
            centis: Math.floor((Date.now() - this.state.startTime) / 10)
        });
    }

    render() {
        return <div onTouchStart={this.keydown} title='Hold "Space" to start.' onTouchEnd={this.keyup}>
            <h1 className={"time " + this.state.timeState}>{format(this.state.centis)}</h1>
        </div>
    }
}

function ResultList(props) {
    const results = props.results;
    return <p className="results">{
        results.map((result) =>
            <ListItem key={result.datetime}
                      value={result.centis}
                      scramble={result.scramble}/>
        )}</p>
}

function ListItem(props) {
    return <span title={props.scramble}>{format(props.value)} </span>;
}

function Statistics(props) {
    const containerStyle = {
        width: (props.stats[100].best !== Infinity)
        + (props.stats[1000].best !== Infinity)
        + 8 + format(props.stats.worst).length + 'em'
    };
    const results = props.results.map((result) => result.centis);
    return <div>
        <div className="stats" style={containerStyle}>
            {results.length > 0 &&
            <div className="columns col-gapless">
                <div className="col-6">Avg:<span className="right">{format(props.stats.avg)}</span></div>
                <div className="col-6 right-stats">Count:<span className="right">{props.stats.n}</span></div>
                <div className="col-6">Best:<span className="right">{format(props.stats.best)}</span></div>
                <div className="col-6 right-stats">Worst:<span className="right">{format(props.stats.worst)}</span>
                </div>
                {avgs.map((avgOf) => [
                    <AvgCurrent key={avgOf} avgOf={avgOf} current={props.stats[avgOf].best} n={props.stats.n}/>,
                    <AvgBest key={avgOf} avgOf={avgOf} best={props.stats[avgOf].best} n={props.stats.n}/>,
                ])}
            </div>
            }
        </div>
        <ResultList results={props.results}/>
    </div>
}

function AvgBest(props) {
    const results = props.results;
    const avgOf = props.avgOf;
    const n = props.n;

    if (props.best === Infinity) {
        return null;
    } else {
        return <div className="col-6 right-stats">Best{avgOf}:<span className="right">{format(props.best)}</span></div>
    }
}

function AvgCurrent(props) {
    const results = props.results;
    const avgOf = props.avgOf;
    const n = props.n;

    if (props.current === Infinity) {
        return null;
    } else {
        return <div className="col-6">Avg{avgOf}:<span className="right left-stats">{format(props.current)}</span></div>
    }
}

class SessionSelector extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selecting: true,
            event: '3x3',
            name: "",
            selectedSession: "",
            sessions: [{name: "main", event: "3x3"}],
            events: ["3x3"]
        };

        this.deleteSession = this.deleteSession.bind(this);
        this.handleSelection = this.handleSelection.bind(this);
        this.handleNameChange = this.handleNameChange.bind(this);
        this.handleEventChange = this.handleEventChange.bind(this);
        this.handleAdding = this.handleAdding.bind(this);
        this.updateActiveSession = this.updateActiveSession.bind(this);
    }

    componentDidMount() {
        if (this.props.auth) {
            axios.get("/api/sessions")
                .then(response => {
                    let sessions = response.data;
                    if (sessions.length === 0) {
                        sessions = [{event: "3x3", name: "main"}];
                    }
                    this.setState({
                        sessions: sessions,
                        selectedSession: sessions[0]
                    });
                    this.updateActiveSession(sessions[0]);
                });
        } else {
            let sessions = JSON.parse(localStorage.getItem("sessions"));
            sessions = sessions === null ? [{event: "3x3", name: "main"}] : sessions;
            this.setState({
                sessions: sessions,
                selectedSession: sessions[0]
            });
            this.updateActiveSession(sessions[0]);
        }
        axios.get('/api/events')
            .then(response => {
                this.setState({
                    events: response.data
                })
            });

    }

    deleteSession() {
        if (confirm('Delete session?')) {
            let deleted = this.state.selectedSession;
            let newSessions = this.state.sessions.filter(value => value !== deleted);
            let updateState = () => {
                if (newSessions.length === 0) {
                    newSessions = [{event: "3x3", name: "main"}]
                }
                const selected = newSessions[0];
                this.setState({
                    sessions: newSessions,
                    selectedSession: selected
                });
                this.updateActiveSession(selected);
            };
            if (this.props.auth) {
                axios.post('/api/deletesessions', {sessions: [deleted]})
                    .then(() => {
                        updateState();
                    })
            } else {
                localStorage.setItem(deleted.event + deleted.name, JSON.stringify([]));
                updateState();
                localStorage.setItem("sessions", JSON.stringify(newSessions));
            }
        }
    }

    updateActiveSession(session) {
        this.setState({
            selectedSession: session
        });
        this.props.setActiveSession(session);
    }

    handleNameChange(e) {
        this.setState({
            name: e.target.value
        });
    }

    handleEventChange(e) {
        this.setState({
            event: e.target.value
        });
    }

    handleAdding(event) {
        event.preventDefault();
        const session = {event: this.state.event, name: this.state.name};
        const newSessions = [...this.state.sessions, session];
        this.setState({
            sessions: newSessions,
            selecting: true,
            name: ""
        });
        this.updateActiveSession(session);
        if (!this.props.auth) {
            localStorage.setItem("sessions", JSON.stringify(newSessions));
        }
    }

    handleSelection(e) {
        if (e.target.value < 0) {
            this.setState({
                selecting: false
            })
        } else {
            this.updateActiveSession(this.state.sessions[e.target.value]);
        }
    }

    render() {
        const chooser = <div>
            <div className="input-group">
                <select className="form-select" onChange={this.handleSelection}
                        value={this.state.sessions.indexOf(this.state.selectedSession)}>
                    {this.state.sessions.map((session, index) => {
                        return <option key={session.event + session.name} value={index}>
                            {session.event + ': ' + session.name}
                        </option>
                    })}
                    <option key={'new'} value={-1}>New session</option>
                </select>
                <button className="btn input-group-btn" onClick={this.deleteSession}>Delete session</button>
            </div>
        </div>;

        const adder =
            <form onSubmit={this.handleAdding}>
                <div className="input-group">
                    <select className="form-select" value={this.state.event} onChange={this.handleEventChange}>
                        {this.state.events.map((event, index) =>
                            <option key={event} value={event}>
                                {event}
                            </option>
                        )}
                    </select>
                    <input className="form-input" type="text" value={this.state.name} placeholder="Session name"
                           onChange={this.handleNameChange}/>
                    <input className="btn btn-primary input-group-btn" type="submit" value="Create"/>
                    <button onClick={() => this.setState({selecting: true})} className="btn input-group-btn">Cancel
                    </button>
                </div>
            </form>;


        return (
            <div className="form-group">
                {this.state.selecting ? chooser : adder}
            </div>
        );
    }
}

class Timer extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            results: [],
            currId: 0,
            stats: calcStats([{centis: 100}, {centis: 100}, {centis: 100}, {centis: 100}, {centis: 100}]),
            scramble: "",
            nextScramble: ""
        };
        this.setActiveSession = this.setActiveSession.bind(this);
        this.saveResults = this.saveResults.bind(this);
        this.saveResult = this.saveResult.bind(this);
        this.updateNextScramble = this.updateNextScramble.bind(this);
        this.updateScramble = this.updateScramble.bind(this);
        this.updateResults = this.updateResults.bind(this);
        this.deleteLast = this.deleteLast.bind(this);
    }

    updateScramble() {
        this.setState((state) => ({
            scramble: state.nextScramble
        }));
        this.updateNextScramble();
    }

    updateNextScramble() {
        axios.get('/api/scramble', {params: {event: this.state.activeSession.event}})
            .then(response => {
                const scramble = response.data;
                this.setState((state) => ({
                    nextScramble: response.data
                }));
            })
    }

    saveResult(res) {
        this.updateScramble();
        const allResults = [...this.state.results, res];
        this.setState((state) => ({
            stats: updateStats(state.stats, allResults),
            scramble: state.nextScramble,
            results: allResults
        }));
        if (this.props.auth) {
            this.saveResults([res]);
        } else {
            localStorage.setItem(this.state.activeSession.event + this.state.activeSession.name, JSON.stringify(allResults));
        }
    }

    saveResults(results) {
        axios.post('/api/saveresults', {
            results: results,
            session: this.state.activeSession
        });
    }

    setActiveSession(session) {
        axios.get('/api/scramble', {params: {event: session.event}})
            .then(response => {
                this.setState({
                    scramble: response.data,
                    nextScramble: response.data
                })
            });
        this.setState({activeSession: session});
        this.updateResults(session);
        setTimeout(this.updateNextScramble, 2000);
    }

    updateResults(session) {
        if (this.props.auth) {
            axios.get('/api/results', {params: {sessionname: session.name, event: session.event}})
                .then(response => {
                    if (response.data != null) {
                        const results = response.data;
                        this.setState({
                            results: results,
                            stats: calcStats(results)
                        });
                    } else {
                        this.setState({
                            results: []
                        });
                    }
                });
        } else {
            let results = JSON.parse(localStorage.getItem(session.event + session.name));
            results = results === null ? [] : results;
            this.setState({
                results: results,
                stats: calcStats(results)
            });
        }
    }

    deleteLast() {
        if (confirm('Delete last solve?')) {
            const newResults = this.state.results.slice(0, -1);
            if (this.props.auth) {
                axios.post('/api/deleteresults', {
                    results: this.state.results.slice(-1),
                    session: this.state.activeSession
                })
                    .then(() => {
                        this.setState({
                            results: newResults,
                            stats: calcStats(newResults)
                        });
                    })
            } else {
                localStorage.setItem(this.state.activeSession.event + this.state.activeSession.name, JSON.stringify(newResults));
                this.setState(state => ({
                    results: newResults,
                    stats: calcStats(newResults)
                }))
            }

        }
    }

    render() {
        return (
            <div className="container columns col-gapless">
                <div className="column col-8 col-lg-12 col-mx-auto">
                    <Header auth={this.props.auth}/>
                    <SessionSelector setActiveSession={this.setActiveSession} auth={this.props.auth}/>
                    <Time saveResult={this.saveResult} scramble={this.state.scramble} centis={this.state.centis}/>
                    <div className="columns">
                        <div className="scramble col-10 col-lg-12 col-mx-auto">{this.state.scramble}</div>
                    </div>
                    <button className="btn" onClick={this.deleteLast}>delete last solve</button>
                    <Statistics results={this.state.results} stats={this.state.stats}/>
                </div>
            </div>
        );
    }
}

class Header extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            profile: {
                name: ""
            }
        }
    }

    componentDidMount() {
        const auth = this.props.auth;
        if (auth) {
            axios.get("/api/userinfo")
                .then(response => {
                    let profile = response.data;
                    if (profile.sub.startsWith("vkontakte")) {
                        profile.formatedName = profile.given_name + ' ' + profile.family_name;
                    } else {//if (profile.sub.startsWith("google-oauth2")) {
                        profile.formatedName = profile.name;
                    }

                    this.setState({
                        profile: profile
                    });
                })
        }
    }

    render() {
        const auth = this.props.auth;

        const loggedIn = <section className="navbar-section">
            <div className="name">{this.state.profile.formatedName}</div>
            <div><a href="/logout" className="btn-link bigfont">Logout</a></div>
        </section>;

        const loggedOut = <section className="navbar-section">
            <a href="/login" className="btn-link bigfont">Login</a>
        </section>;

        return (<div>

                <header className="navbar timernavbar bigfont">
                    <section className="navbar-section title">T-timer</section>
                    {auth ? loggedIn : loggedOut}
                </header>
            </div>
        )
    }
}

function currentDatetime() {
    let d = new Date();
    let s = d.toISOString();
    return s.replace('T', ' ').slice(0, -2);
}

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            app: <div/>
        }
    }

    componentDidMount() {
        axios.get('/api/isauthenticated')
            .then(response => {
                const auth = response.data;
                this.setState({
                    app: <Timer auth={auth}/>
                })
            });
    }

    render() {
        return this.state.app;
    }
}

ReactDOM.render(
    <App/>,
    document.getElementById('app')
);