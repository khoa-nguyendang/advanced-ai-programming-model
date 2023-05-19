import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import { IconButton, Modal, Toolbar, Tooltip, Typography } from '@mui/material';
import Paper from '@mui/material/Paper';
import { styled } from '@mui/material/styles';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import * as React from 'react';
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useEffectOnce } from 'react-use';
import AddUser from './AddUser';

const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
        border: 0,
    },
}));

function createData(name, calories, fat, carbs, protein) {
    return { name, calories, fat, carbs, protein };
}

const rows = [
    createData('Frozen yoghurt', 159, 6.0, 24, 4.0),
    createData('Ice cream sandwich', 237, 9.0, 37, 4.3),
    createData('Eclair', 262, 16.0, 24, 6.0),
    createData('Cupcake', 305, 3.7, 67, 4.3),
    createData('Gingerbread', 356, 16.0, 49, 3.9),
];

export default function Users() {
    const [users, setUsers] = React.useState([]);
    const [open, setOpen] = React.useState(false);
    const handleOpen = () => setOpen(true);
    const handleClose = () => {
        setOpen(false);
        toast.success("User Added sucessfully!", {
            position: "bottom-right",
            autoClose: 5000,
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
            progress: undefined,
            theme: "colored",
            });
    };

    const fetchUsers = async function () {
        console.log("calling the api fetchUsers");
        let body = {
            userRoleId: "1",
            userGroupIds: [],
            userIds: [],
            referenceIds: [],
            currentPage: "0",
            pageSize: "100"
        }
        await fetch('http://localhost:8081/api/v1/user/get-user', {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {
                'Content-type': 'application/json; charset=UTF-8',
            },
        })
            .then((response) => response.json())
            .then((data) => {
                setUsers(data.users);
            })
            .catch((err) => {
                console.log(err.message);
            });
    }

    const handleDelete = function (userId) {
        console.log("calling the api handleDelete");
        let body = {
            userIds: [userId],
            referenceIds: [],
        };
        fetch('http://localhost:8081/api/v1/user/delete', {
            method: 'POST',
            body: JSON.stringify(body),
            headers: {
                'Content-type': 'application/json; charset=UTF-8',
            },
        })
            .then((response) => response.json())
            .then((data) => {
                toast.success("User Deleted sucessfully!", {
                    position: "bottom-right",
                    autoClose: 5000,
                    hideProgressBar: false,
                    closeOnClick: true,
                    pauseOnHover: true,
                    draggable: true,
                    progress: undefined,
                    theme: "colored",
                    });
            })
            .catch((err) => {
                console.log(err.message);
            });
    }

    useEffectOnce(() => {
        fetchUsers();
        return () => {
            console.log('Running clean-up of effect on unmount')
        }
    });
    console.log("users: ", users);

    return (
        <Paper >
            <Toolbar sx={{ pl: { sm: 2 }, pr: { xs: 1, sm: 1 } }} >
                <Typography
                    sx={{ flex: '1 1 100%' }}
                    variant="h6"
                    id="tableTitle"
                    component="div"
                >
                    Users
                </Typography>

                <Tooltip title="Add">
                    <IconButton onClick={handleOpen}>
                        <AddIcon />
                    </IconButton>
                </Tooltip>
            </Toolbar>
            <TableContainer component={Paper}>
                <Table sx={{ minWidth: 700 }} aria-label="customized table">
                    <TableHead>
                        <TableRow>
                            <StyledTableCell align="left">User Id</StyledTableCell>
                            <StyledTableCell align="left">User Name</StyledTableCell>
                            <StyledTableCell align="left">Reference Id</StyledTableCell>
                            <StyledTableCell align="left">Full Name</StyledTableCell>
                            <StyledTableCell align="left">Email</StyledTableCell>
                            <StyledTableCell align="left">Images</StyledTableCell>
                            <StyledTableCell align="right">Actions</StyledTableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {!Array.isArray(users) ? null : users?.map((row) => {
                            let userInfo = JSON.parse(row.userInfo);
                            return (
                                <StyledTableRow key={row.userId}>
                                    <StyledTableCell component="th" scope="row" align="left">
                                        {row.userId}
                                    </StyledTableCell>
                                    <StyledTableCell align="left">{row.userName}</StyledTableCell>
                                    <StyledTableCell align="left">{row.referenceId}</StyledTableCell>
                                    <StyledTableCell align="left">{userInfo.first + ' ' + userInfo.last}</StyledTableCell>
                                    <StyledTableCell align="left">{userInfo.email}</StyledTableCell>
                                    <StyledTableCell align="left">
                                        {userInfo.email}
                                    </StyledTableCell>
                                    <StyledTableCell align="right">
                                        <EditIcon></EditIcon>
                                        <DeleteIcon onClick={() => handleDelete(row.userId)}></DeleteIcon>
                                    </StyledTableCell>
                                </StyledTableRow>
                            )
                        })}
                    </TableBody>
                </Table>
            </TableContainer>
            <ToastContainer />
            <Modal
                open={open}
                onClose={handleClose}
                aria-labelledby="modal-modal-title"
                aria-describedby="modal-modal-description"
            >
                <AddUser onComplete={handleClose}/>
            </Modal>
        </Paper>

    );
}
