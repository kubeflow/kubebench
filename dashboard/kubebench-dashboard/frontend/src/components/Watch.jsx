import React from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import ScheduleIcon from '@material-ui/icons/Schedule';
import HighlightOffIcon from '@material-ui/icons/HighlightOff';
import DoneIcon from '@material-ui/icons/Done';
import LinearProgress from '@material-ui/core/LinearProgress';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';

import { openModal, closeModal, fetchJobs, closeSnack } from "../actions";

import JobInfo from './JobInfo'

const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
    },
    running: {
        color: '#8b8ffb',
    },
    failed: {
        color: '#f26363',
    },
    finished: {
        color: '#63f291',
    },
    progress: {
        height: 10,
        margin: 10,
    },
    close: {
        padding: theme.spacing.unit / 2,
    },
});


class Watch extends React.Component {
    componentDidMount() {
        this.props.fetchJobs()
    }

    open = (id) => event => {
        this.props.openModal(id);
    }

    deleteJob = (name) => {
        this.props.deleteJob(name)
    }

    render () {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <h1>Watch</h1>
                <hr />
                {this.props.loading && <LinearProgress className={classes.progress}/>}
                <List component="nav">
                    {this.props.jobs.map((job, i) => {
                        let icon;
                        if (job.status === 'Running') {
                            icon = (<ScheduleIcon className={classes.running}/>)
                        } else if (job.status === 'Failed') {
                            icon = (<HighlightOffIcon className={classes.failed}/>)
                        } else {
                            icon = (<DoneIcon className={classes.finished}/>)
                        }
                        return (
                            <ListItem button key={i} onClick={this.open(i)}>
                                <ListItemIcon>
                                    {icon}
                                </ListItemIcon>
                                <ListItemText inset primary={job.name} />
                            </ListItem>
                        );
                    })}
                </List>
                <JobInfo
                    close={this.props.closeModal}
                    open={this.props.modalOpen}
                    name={this.props.name}
                />
                 <Snackbar
                    anchorOrigin={{
                        vertical: 'bottom',
                        horizontal: 'center',
                    }}
                    open={this.props.snackOpen}
                    autoHideDuration={6000}
                    onClose={this.props.closeSnack}
                    message={<span id="message-id">{this.props.snackText}</span>}
                    action={[
                        <IconButton
                            key="close"
                            aria-label="Close"
                            color="inherit"
                            className={classes.close}
                            onClick={this.props.closeSnack}
                        >
                            <CloseIcon />
                        </IconButton>,
                    ]}
                />
            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        jobs: state.jobsList,
        modalOpen: state.modalOpen,
        currentId: state.current,
        name: state.currentName,
        loading: state.loading,
        snackOpen: state.snackOpen,
        snackText: state.snackText
    };
};

Watch.propTypes = {
    classes: PropTypes.object.isRequired,
};


export default connect(mapStateToProps, { openModal, closeModal, fetchJobs, closeSnack })(withStyles(styles)(Watch));