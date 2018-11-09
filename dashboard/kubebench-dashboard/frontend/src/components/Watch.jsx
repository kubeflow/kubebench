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
import TextField from '@material-ui/core/TextField';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import Button from '@material-ui/core/Button';


import { openModal, closeModal, fetchJobs, closeSnack, filterJobs, changeType } from "../actions";

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
    textField: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
    },
    filter: {
        margin: '0 auto',
        textAlign: 'center',
    },
    button: {
        margin: theme.spacing.unit / 2,
        padding: theme.spacing.unit / 2,
    }
});


class Watch extends React.Component {
    componentDidMount() {
        this.props.fetchJobs();
    }

    onFilter = event => {
        this.props.filterJobs(event.target.value);
    }

    open = (id) => event => {
        this.props.openModal(id);
    }

    deleteJob = (name) => {
        this.props.deleteJob(name);
    }

    update = (event) => {
        this.props.fetchJobs();
    }

    handleType = (name) => event => {
        this.props.changeType(name, event.target.checked);
    }

    render () {
        const { classes } = this.props;
        return (
            <div className={classes.root}>
                <h1>Monitor</h1>
                <hr />
                {this.props.loading && <LinearProgress className={classes.progress}/>}
                <div className={classes.filter}>
                    <FormGroup row>
                        <TextField
                            id="outlined-name"
                            label="Name"
                            className={classes.textField}
                            value={this.props.filter}
                            onChange={this.onFilter}
                            margin="normal"
                            variant="outlined"
                        />
                        {
                                Object.keys(this.props.filterType).map((filter, i) => {
                                return(
                                    <FormControlLabel
                                        key={i}
                                        control={
                                            <Switch
                                            checked={this.props.filterType[filter]}
                                            onChange={this.handleType(filter)}
                                            value={filter}
                                            color={"secondary"}
                                            />
                                        }
                                        label={filter}
                                    />
                                );
                            })
                        }
                    </FormGroup>
                    <Button color={"secondary"} onClick={this.update} variant={"raised"}>
                        Update
                    </Button>
                </div>
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
        jobs: state.filteredJobsList,
        modalOpen: state.modalOpen,
        currentId: state.current,
        name: state.currentName,
        loading: state.loading,
        snackOpen: state.snackOpen,
        snackText: state.snackText,
        filter: state.filter,
        filterType: state.filterType,
    };
};

Watch.propTypes = {
    classes: PropTypes.object.isRequired,
};


export default connect(mapStateToProps, { openModal, closeModal, fetchJobs, closeSnack, filterJobs, changeType })(withStyles(styles)(Watch));