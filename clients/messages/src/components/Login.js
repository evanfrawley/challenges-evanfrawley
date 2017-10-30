import React from 'react';
import {Row, Input} from 'react-materialize';

class Login extends React.Component {
    render() {
        return(
            <div className={"LoginContainer"}>
                <Row>
                    <Input type="password" label="password" s={12} />
                    <Input type="email" label="Email" s={12} />
                </Row>
                <div>
                  <p>Don't have an account? <a href={"/signup"}>Sign up now!</a></p>
                </div>
            </div>
        )
    }
}

export default Login