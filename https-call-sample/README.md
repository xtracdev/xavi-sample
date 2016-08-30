## HTTPs Call Using Load Balancer

This sample illustrates how to make an HTTPs call using a load
balancer.

This will not run on your machine unless you do the following:

1. In your golang root (try /usr/local/go) locate src/crypto/tls/generate_cert.go,
and build it.
2. Use generate_cert to generate a cacert for this example:
<pre>
generate_cert -ca -host `hostname`
</pre>
3. Copy the generated cert.pem into this directory, and insert the contents
of cert.pem and key.pem into inposter-https.json.

To run the sample, start (or restart) mountebank, run `sample-setup.sh`,
then 

<pre>
env XAVI_KVSTORE_URL=file:///`pwd`/config go run lbsample.go
</pre>


