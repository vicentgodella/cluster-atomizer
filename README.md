Cluster atomizer is a tool to do immutable deployments on Amazon Web Services.

# cluster command

## list clusters

Usage:
`cloud-atom cluster [-p <profile>] [-r <region>] list [-a <application>] [-j <project>]`

For example:

`cloud-atom cluster -p my-aws-profile  -j my_project list`

```
+-----------------------------+-----------------+----------------------------------------------------------------+
|            NAME             | INSTANCES COUNT |                      LOAD BALANCERS LIST                       |
+-----------------------------+-----------------+----------------------------------------------------------------+
| my_project-http-v001        |               0 | my-project-int-http,my-project-pub-http                        |
| my_project-workers-v003     |               0 | my-project-int-workers                                         |
+-----------------------------+-----------------+----------------------------------------------------------------+
```
