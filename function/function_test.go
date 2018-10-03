package function

/**********************************************************************/
/*  The use of reflection breaks the assurances garaneed by golangs   */
/*  type system. Thus it should never be used.                        */
/*  We use refelction anyway, thus we have to be extra though in      */
/*  for testing of type errors.                                       */
/*  When RunFunc is provided with a function with signature that      */
/*  violates the specification RunFunc should panic                   */
/*  RunFunc should never panic in runtime due to reflection           */
/*  operations if its provided with a corect function                 */
/*                                                                    */
/*  TODO: test if RunFunc acording to above specification             */
/*      *) https://github.com/onsi/gomega                             */
/*         Has nice matchesrs for testing panic                       */
/*                                                                    */
/**********************************************************************/

import (
	"testing"
)

func TestRuning(t *testing.T) {

}
