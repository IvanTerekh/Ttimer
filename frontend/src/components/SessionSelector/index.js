import React from 'react';

const axios = require('axios');

export default class SessionSelector extends React.Component {
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
            axios.get('/api/sessions')
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
        if (window.confirm('Delete session?')) {
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
