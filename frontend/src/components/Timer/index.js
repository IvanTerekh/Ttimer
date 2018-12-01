import React from 'react';
import Time from '../Time';
import Statistics from '../Statistics';
import SessionSelector from '../SessionSelector';
import {calcStats, updateStats} from '../../logic/stats';
import Header from '../Header';
import './Timer.css';

const axios = require('axios');

export default class Timer extends React.Component {
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
                    nextScramble: scramble
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
        if (window.confirm('Delete last solve?')) {
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
