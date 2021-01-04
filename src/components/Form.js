import React, { useState } from 'react';
import axios from 'axios';

function Form() {

    const [firstWord, setFirstWord] = useState('');
    const [secondWord, setSecondWord] = useState('');
    const [results, setResults] = useState(null);
    const [isError, setIsError] = useState(false);
    const [isLoading, setIsLoading] = useState(false);


    const submitForm = async (e) => {
        e.preventDefault()
        let body = {
            firstWord,
            secondWord
        };
        
        setIsLoading(true);

        axios
          .post("http://localhost:8000/api/postWords", body)
          .then(function(response) {
            setResults(response.data)
          })
          .catch(function(error) {
            setIsError(true)
          });

        setIsLoading(false)
    }

  return (
    <div>
        <form onSubmit={(e) => {submitForm(e)}}>
        <label>
          First Word:
          <input type="text" name="firstWord" 
            onKeyUp={(e) => setFirstWord(e.target.value)}
            />
        </label>
        <label>
          Second Word:
          <input type="text" name="secondWord" 
            onKeyUp={(e) => setSecondWord(e.target.value)}
          />
        </label>

        <button type="submit">Get distances!</button>   

        { isError &&  <div>Something went wrong ...</div > }

        { isLoading && <div>Loading ...</div> }

        { results && !isLoading &&
            (
                <table>
                    <tr>
                        <td>the first word</td><td>{firstWord}</td>
                    </tr>
                    <tr>
                        <td>the second word</td> <td>{secondWord}</td>
                    </tr>
                    <tr>
                        <td>the absolute length difference</td> <td>{results.absoluteDifference}</td>
                    </tr>
                    <tr>
                        <td>the levenstein difference</td> <td>{results.levensteinDifference}</td>
                    </tr>
                    <tr>
                        <td>third party leveinstein difference</td><td>{results.thirdPartylevensteinDifference}</td>
                    </tr>
                </table>
            )
        }
        </form>
    </div>
  );
}

export default Form;