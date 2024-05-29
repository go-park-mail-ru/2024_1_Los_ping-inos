# Details

Date : 2024-05-29 20:00:56

Directory /home/longwave/techpark/workisdone/2024_1_Los_ping-inos/internal

Total : 62 files,  12069 codes, 0 comments, 1871 blanks, all 13940 lines

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)

## Files
| filename | language | code | comment | blank | total |
| :--- | :--- | ---: | ---: | ---: | ---: |
| [internal/auth/auth_easyjson.go](/internal/auth/auth_easyjson.go) | Go | 1,337 | 0 | 56 | 1,393 |
| [internal/auth/auth_models.go](/internal/auth/auth_models.go) | Go | 157 | 0 | 25 | 182 |
| [internal/auth/delivery/grpc.go](/internal/auth/delivery/grpc.go) | Go | 58 | 0 | 10 | 68 |
| [internal/auth/delivery/http.go](/internal/auth/delivery/http.go) | Go | 482 | 0 | 67 | 549 |
| [internal/auth/delivery/http_test.go](/internal/auth/delivery/http_test.go) | Go | 172 | 0 | 31 | 203 |
| [internal/auth/interfaces.go](/internal/auth/interfaces.go) | Go | 44 | 0 | 6 | 50 |
| [internal/auth/mocks/core_mocks.go](/internal/auth/mocks/core_mocks.go) | Go | 284 | 0 | 52 | 336 |
| [internal/auth/proto/auth.pb.go](/internal/auth/proto/auth.pb.go) | Go | 402 | 0 | 55 | 457 |
| [internal/auth/proto/auth.proto](/internal/auth/proto/auth.proto) | Protocol Buffers | 27 | 0 | 8 | 35 |
| [internal/auth/proto/auth_grpc.pb.go](/internal/auth/proto/auth_grpc.pb.go) | Go | 129 | 0 | 18 | 147 |
| [internal/auth/repo/interest.go](/internal/auth/repo/interest.go) | Go | 139 | 0 | 27 | 166 |
| [internal/auth/repo/payment.go](/internal/auth/repo/payment.go) | Go | 69 | 0 | 10 | 79 |
| [internal/auth/repo/person.go](/internal/auth/repo/person.go) | Go | 201 | 0 | 35 | 236 |
| [internal/auth/repo/repo_test.go](/internal/auth/repo/repo_test.go) | Go | 686 | 0 | 189 | 875 |
| [internal/auth/repo/session.go](/internal/auth/repo/session.go) | Go | 48 | 0 | 10 | 58 |
| [internal/auth/usecase/profile.go](/internal/auth/usecase/profile.go) | Go | 138 | 0 | 29 | 167 |
| [internal/auth/usecase/profile_test.go](/internal/auth/usecase/profile_test.go) | Go | 606 | 0 | 77 | 683 |
| [internal/auth/usecase/usecase.go](/internal/auth/usecase/usecase.go) | Go | 177 | 0 | 34 | 211 |
| [internal/auth/usecase/usecase_test.go](/internal/auth/usecase/usecase_test.go) | Go | 647 | 0 | 94 | 741 |
| [internal/cmd/auth/main.go](/internal/cmd/auth/main.go) | Go | 157 | 0 | 36 | 193 |
| [internal/cmd/feed/main.go](/internal/cmd/feed/main.go) | Go | 132 | 0 | 33 | 165 |
| [internal/cmd/images/main.go](/internal/cmd/images/main.go) | Go | 66 | 0 | 12 | 78 |
| [internal/feed/delivery/http.go](/internal/feed/delivery/http.go) | Go | 316 | 0 | 41 | 357 |
| [internal/feed/feed_easyjson.go](/internal/feed/feed_easyjson.go) | Go | 1,269 | 0 | 45 | 1,314 |
| [internal/feed/feed_models.go](/internal/feed/feed_models.go) | Go | 150 | 0 | 19 | 169 |
| [internal/feed/interfaces.go](/internal/feed/interfaces.go) | Go | 37 | 0 | 6 | 43 |
| [internal/feed/mocks/core_mocks.go](/internal/feed/mocks/core_mocks.go) | Go | 263 | 0 | 46 | 309 |
| [internal/feed/repo/chat.go](/internal/feed/repo/chat.go) | Go | 136 | 0 | 21 | 157 |
| [internal/feed/repo/claim.go](/internal/feed/repo/claim.go) | Go | 52 | 0 | 8 | 60 |
| [internal/feed/repo/interest.go](/internal/feed/repo/interest.go) | Go | 92 | 0 | 19 | 111 |
| [internal/feed/repo/like.go](/internal/feed/repo/like.go) | Go | 118 | 0 | 23 | 141 |
| [internal/feed/repo/person.go](/internal/feed/repo/person.go) | Go | 80 | 0 | 17 | 97 |
| [internal/feed/repo/postgres.go](/internal/feed/repo/postgres.go) | Go | 10 | 0 | 4 | 14 |
| [internal/feed/repo/repo_test.go](/internal/feed/repo/repo_test.go) | Go | 911 | 0 | 222 | 1,133 |
| [internal/feed/usecase/chats.go](/internal/feed/usecase/chats.go) | Go | 31 | 0 | 6 | 37 |
| [internal/feed/usecase/chats_test.go](/internal/feed/usecase/chats_test.go) | Go | 145 | 0 | 26 | 171 |
| [internal/feed/usecase/usecase.go](/internal/feed/usecase/usecase.go) | Go | 107 | 0 | 22 | 129 |
| [internal/feed/usecase/usecase_test.go](/internal/feed/usecase/usecase_test.go) | Go | 470 | 0 | 68 | 538 |
| [internal/image/delivery/grpc/grpc.go](/internal/image/delivery/grpc/grpc.go) | Go | 32 | 0 | 6 | 38 |
| [internal/image/delivery/http/http.go](/internal/image/delivery/http/http.go) | Go | 191 | 0 | 43 | 234 |
| [internal/image/delivery/http/http_test.go](/internal/image/delivery/http/http_test.go) | Go | 26 | 0 | 9 | 35 |
| [internal/image/image_easyjson.go](/internal/image/image_easyjson.go) | Go | 77 | 0 | 9 | 86 |
| [internal/image/image_models.go](/internal/image/image_models.go) | Go | 30 | 0 | 4 | 34 |
| [internal/image/interfaces.go](/internal/image/interfaces.go) | Go | 17 | 0 | 4 | 21 |
| [internal/image/mocks/core_mocks.go](/internal/image/mocks/core_mocks.go) | Go | 76 | 0 | 16 | 92 |
| [internal/image/protos/gen/image.pb.go](/internal/image/protos/gen/image.pb.go) | Go | 202 | 0 | 28 | 230 |
| [internal/image/protos/gen/image_grpc.pb.go](/internal/image/protos/gen/image_grpc.pb.go) | Go | 94 | 0 | 16 | 110 |
| [internal/image/protos/proto/image.proto](/internal/image/protos/proto/image.proto) | Protocol Buffers | 15 | 0 | 7 | 22 |
| [internal/image/repo/image.go](/internal/image/repo/image.go) | Go | 175 | 0 | 38 | 213 |
| [internal/image/repo/repo_test.go](/internal/image/repo/repo_test.go) | Go | 179 | 0 | 47 | 226 |
| [internal/image/usecase/image.go](/internal/image/usecase/image.go) | Go | 69 | 0 | 20 | 89 |
| [internal/image/usecase/image_test.go](/internal/image/usecase/image_test.go) | Go | 174 | 0 | 28 | 202 |
| [internal/logs/log.go](/internal/logs/log.go) | Go | 21 | 0 | 6 | 27 |
| [internal/pkg/middlewares.go](/internal/pkg/middlewares.go) | Go | 96 | 0 | 14 | 110 |
| [internal/pkg/pkg_test.go](/internal/pkg/pkg_test.go) | Go | 73 | 0 | 15 | 88 |
| [internal/pkg/requests.go](/internal/pkg/requests.go) | Go | 1 | 0 | 1 | 2 |
| [internal/pkg/responce_test.go](/internal/pkg/responce_test.go) | Go | 33 | 0 | 14 | 47 |
| [internal/pkg/responces.go](/internal/pkg/responces.go) | Go | 60 | 0 | 14 | 74 |
| [internal/pkg/security.go](/internal/pkg/security.go) | Go | 40 | 0 | 8 | 48 |
| [internal/types/errors.go](/internal/types/errors.go) | Go | 10 | 0 | 5 | 15 |
| [internal/types/types.go](/internal/types/types.go) | Go | 8 | 0 | 4 | 12 |
| [internal/types/types_test.go](/internal/types/types_test.go) | Go | 25 | 0 | 8 | 33 |

[Summary](results.md) / Details / [Diff Summary](diff.md) / [Diff Details](diff-details.md)